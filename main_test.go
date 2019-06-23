package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const expectedHash = "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
const knownPassword = "angryMonkey"

var _ = Describe("Main", func() {
    var (
       returnedHash string
       testPass string
    )

    Context("#hashString", func() {
        JustBeforeEach(func () {
           returnedHash = hashString(testPass)
        })

        Context("with a known password", func() {
            BeforeEach(func() {
               testPass = knownPassword
            })

            It("should return the known hash", func() {
                Expect(returnedHash).To(Equal(expectedHash))
            })
        })
    })
})
