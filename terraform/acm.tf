# Generate a self-signed certificate using TLS provider
resource "tls_private_key" "main" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "main" {
  private_key_pem = tls_private_key.main.private_key_pem

  subject {
    common_name  = "${var.project_name}-${var.environment}.localhost"
    organization = var.project_name
  }

  # Certificate is valid for 1 year
  validity_period_hours = 8760

  # Allowed uses
  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]

  # Add the ALB DNS name as a DNS SAN after it's created
  dns_names = [
    "${var.project_name}-${var.environment}.localhost",
    "localhost"
  ]
}

# Import the self-signed certificate to ACM
resource "aws_acm_certificate" "main" {
  private_key      = tls_private_key.main.private_key_pem
  certificate_body = tls_self_signed_cert.main.cert_pem

  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Name        = "${var.project_name}-${var.environment}-cert"
    Environment = var.environment
  }
}
