// Declaramos el paquete utils para utilidades de JWT.
package utils

import (
	// Importamos fmt para formatear strings.
	"fmt"
	"os"

	// Importamos time para manejar tiempos de expiración.
	"time"

	// Importamos las librerías de JWT.
	"github.com/golang-jwt/jwt/v5"
)

// JWTSecret es la clave secreta para firmar tokens JWT.
var JWTSecret = []byte("multicatalogo-secret-key-change-in-production")

// ConfigureJWTSecret lee JWT_SECRET. Fuera de APP_ENV=dev es obligatorio y debe
// tener al menos 32 caracteres; asi nunca se usa la clave por defecto publica.
func ConfigureJWTSecret() error {
	secret := os.Getenv("JWT_SECRET")
	if os.Getenv("APP_ENV") != "dev" {
		if len(secret) < 32 {
			return fmt.Errorf("JWT_SECRET es obligatorio (minimo 32 caracteres) cuando APP_ENV no es dev")
		}
	}
	if secret != "" {
		JWTSecret = []byte(secret)
	}
	return nil
}

// Claims representa los claims del token JWT.
type Claims struct {
	Email string `json:"email"`
	Rol   string `json:"rol"`
	jwt.RegisteredClaims
}

// GenerateToken genera un token JWT para un usuario.
func GenerateToken(email string, rol string) (string, error) {
	// Creamos los claims con expiración de 2 horas.
	expirationTime := time.Now().Add(2 * time.Hour)
	claims := &Claims{
		Email: email,
		Rol:   rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "multicatalogo",
		},
	}

	// Creamos el token con los claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmamos el token con la clave secreta.
	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", fmt.Errorf("error al firmar token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken valida un token JWT y retorna los claims.
func ValidateToken(tokenString string) (*Claims, error) {
	// Parseamos el token.
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validamos que el método de firma sea el correcto.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inválido: %v", token.Header["alg"])
		}
		return JWTSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error al parsear token: %w", err)
	}

	// Extraemos los claims.
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("token inválido")
}
