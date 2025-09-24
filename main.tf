terraform {
  required_providers {
    dob = {
      source = "liatr.io/terraform/devops-bootcamp"
    }
  }
}

provider "dob" {
  endpoint = "http://localhost:8080"
}

resource "dob_engineer" "Madi" {
    name = "Madi"
    email = "madi@liatrio.com"
}

resource "dob_engineer" "Colin" {
    name = "Colin"
    email = "colin@liatrio.com"
}

resource "dob_engineer" "Angel" {
    name = "Angel"
    email = "angel@liatrio.com"
}

resource "dob_engineer" "Austin" {
    name = "Austin"
    email = "austin@liatrio.com"
}

resource "dob_engineer" "Jack" {
    name = "Jack"
    email = "jack@liatrio.com"
}

resource "dob_dev" "example" {
    name = "Dev Team #1"
    engineers = [
      dob_engineer.Angel.id, 
      dob_engineer.Colin.id
    ]
}

resource "dob_ops" "example" {
    name = "Ops Team #1"
    engineers = [
      dob_engineer.Madi.id, 
      dob_engineer.Austin.id, 
      dob_engineer.Jack.id
    ]
}
