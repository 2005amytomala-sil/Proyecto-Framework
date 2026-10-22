package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"multicatalogo-backend/interfaces"
	"multicatalogo-backend/models"
	"multicatalogo-backend/utils"
	"os"
	"strings"
)

type AuthController struct{ repo interfaces.AuthRepository }

func NewAuthController(repo interfaces.AuthRepository) *AuthController {
	return &AuthController{repo: repo}
}
func apiError(c *fiber.Ctx, status int, message string, details ...interface{}) error {
	var detail interface{}
	if len(details) > 0 && os.Getenv("APP_ENV") == "dev" {
		detail = details[0]
	}
	return c.Status(status).JSON(models.APIError{Status: status, Message: message, Details: detail})
}

func (ac *AuthController) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return apiError(c, 400, "Cuerpo de petición inválido")
	}
	req.Nombre = strings.TrimSpace(req.Nombre)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Ref = strings.ToUpper(strings.TrimSpace(req.Ref))
	if req.Nombre == "" || req.Email == "" || req.Password == "" {
		return apiError(c, 400, "Nombre, email y contraseña son obligatorios")
	}
	if !strings.Contains(req.Email, "@") {
		return apiError(c, 400, "Email inválido")
	}
	if len(req.Password) < 6 {
		return apiError(c, 400, "La contraseña debe tener al menos 6 caracteres")
	}
	if ac.repo == nil {
		return apiError(c, 503, "Base de datos no disponible")
	}
	if _, err := ac.repo.BuscarUsuarioPorEmail(req.Email); err == nil {
		return apiError(c, 409, "El email ya está registrado")
	} else if !errors.Is(err, interfaces.ErrNotFound) {
		return apiError(c, 500, "Error al validar el registro")
	}
	var parentID *int
	nivel := 0
	if req.Ref != "" {
		id, err := ac.repo.BuscarReferidoPorCodigo(req.Ref)
		if err != nil {
			if errors.Is(err, interfaces.ErrNotFound) {
				return apiError(c, 400, "Código de referido inválido")
			}
			return apiError(c, 500, "Error al validar el código")
		}
		parentID = &id
		nivel, err = ac.repo.ObtenerNivelReferido(id)
		if err != nil {
			return apiError(c, 500, "Error al consultar el referente")
		}
		nivel++
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return apiError(c, 500, "Error al procesar la contraseña")
	}
	tx, err := ac.repo.Begin()
	if err != nil {
		return apiError(c, 500, "Error al iniciar el registro")
	}
	defer tx.Rollback()
	if err = ac.repo.CrearUsuarioTx(tx, req.Email, string(hash), req.Nombre, "cliente"); err != nil {
		return apiError(c, 409, "No se pudo crear la cuenta")
	}
	var b [4]byte
	if _, err = rand.Read(b[:]); err != nil {
		return apiError(c, 500, "Error al generar el código")
	}
	codigo := strings.ToUpper(hex.EncodeToString(b[:]))
	nodoID, err := ac.repo.CrearNodoReferidoTx(tx, req.Email, req.Nombre, codigo, nivel, parentID)
	if err != nil {
		return apiError(c, 500, "Error al crear el nodo de red")
	}
	if err = ac.repo.VincularUsuarioReferidoTx(tx, req.Email, nodoID, parentID); err != nil {
		return apiError(c, 500, "Error al vincular la cuenta")
	}
	if err = tx.Commit(); err != nil {
		return apiError(c, 500, "Error al confirmar el registro")
	}
	return c.Status(201).JSON(fiber.Map{"message": "Usuario registrado exitosamente", "email": req.Email, "codigo": codigo})
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return apiError(c, 400, "Cuerpo de petición inválido")
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Email == "" || req.Password == "" {
		return apiError(c, 400, "Email y contraseña son obligatorios")
	}
	if ac.repo == nil {
		return apiError(c, 503, "Base de datos no disponible")
	}
	u, err := ac.repo.BuscarUsuarioPorEmail(req.Email)
	if err != nil {
		if errors.Is(err, interfaces.ErrNotFound) {
			return apiError(c, 401, "Credenciales incorrectas")
		}
		return apiError(c, 500, "No se pudo validar el acceso")
	}
	if !u.Activo {
		return apiError(c, 403, "Cuenta desactivada")
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
		return apiError(c, 401, "Credenciales incorrectas")
	}
	token, err := utils.GenerateToken(u.Email, u.Rol)
	if err != nil {
		return apiError(c, 500, "Error al generar token")
	}
	return c.JSON(fiber.Map{"token": token, "email": u.Email, "rol": u.Rol})
}
