package collision

import (
	"github.com/teomat/mater/vect"
	"log"
)

func clamp(val, min, max float64) float64 {
	if val < min {
		return min
	} else if val > max {
		return max
	}
	return val
}

type hashValue uint64

const hashCoef = hashValue(3344921057)

func hashPair(a, b hashValue) hashValue {
	return (a * hashCoef) ^ (b * hashCoef)
}

func k_scalar_body(body *Body, r, n vect.Vect) float64 {
	rcn := vect.Cross(r, n)
	return body.invMass + body.invI*rcn*rcn
}

func k_scalar(a, b *Body, r1, r2, n vect.Vect) float64 {
	value := k_scalar_body(a, r1, n) + k_scalar_body(b, r2, n)
	if value == 0.0 {
		log.Printf("Warning: Unsolvable collision or constraint.")
	}
	return value
}

func relative_velocity(a, b *Body, r1, r2 vect.Vect) vect.Vect {
	v1_sum := vect.Add(a.Velocity, vect.Mult(vect.Perp(r1), a.AngularVelocity))
	v2_sum := vect.Add(b.Velocity, vect.Mult(vect.Perp(r2), b.AngularVelocity))

	return vect.Sub(v2_sum, v1_sum)
}

func normal_relative_velocity(a, b *Body, r1, r2, n vect.Vect) float64 {
	return vect.Dot(relative_velocity(a, b, r1, r2), n)
}

func apply_impulses(a, b *Body, r1, r2, j vect.Vect) {
	apply_impulse(a, vect.Mult(j, -1), r1)
	apply_impulse(b, j, r2)
}

func apply_impulse(body *Body, j, r vect.Vect) {
	body.Velocity.Add(vect.Mult(j, body.invMass))
	body.AngularVelocity += body.invI * vect.Cross(r, j)
}

func apply_bias_impulses(a, b *Body, r1, r2, j vect.Vect) {
	apply_bias_impulse(a, vect.Mult(j, -1), r1)
	apply_bias_impulse(b, j, r2)
}

func apply_bias_impulse(body *Body, j, r vect.Vect) {
	body.v_bias.Add(vect.Mult(j, body.invMass))
	body.w_bias += body.invI * vect.Cross(r, j)
}
