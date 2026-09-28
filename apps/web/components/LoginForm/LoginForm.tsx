"use client"

import { Button } from "@/components/ui/button"
import { Field, FieldDescription, FieldError } from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { requestMagicLink, verifyLoginCode } from "@/lib/auth"
import { useForm } from "@tanstack/react-form"
import { useMutation } from "@tanstack/react-query"
import z from "zod"
import { useState } from "react"

const formSchema = z.object({
    email: z.email()
})

const codeSchema = z.object({
    code: z.string().regex(/^\d{6}$/, "Enter the 6-digit code")
})

export function LoginForm() {
    const [sentTo, setSentTo] = useState<string | null>(null)

    const sendLink = useMutation({
        mutationFn: (email: string) => requestMagicLink(email),
        onSuccess: (_, email) => setSentTo(email)
    })

    const { handleSubmit, Field: FromField } = useForm({
        defaultValues: {
            email: ""
        },
        validators: {
            onSubmit: formSchema
        },
        onSubmit: ({ value }) => {
            sendLink.mutate(value.email)
        }
    })

    if (sentTo) {
        return <LoginCodeForm email={sentTo} onUseDifferentEmail={() => {
            sendLink.reset()
            setSentTo(null)
        }} />
    }

    return <form
        className="flex flex-col gap-4 w-full"
        id="login-form"
        onSubmit={(e) => {
            e.preventDefault()
            handleSubmit()
        }}
    >
        <FromField name="email">
            {(field) => {
                const isInvalid =
                    field.state.meta.isTouched && !field.state.meta.isValid
                return (
                    <Field data-invalid={isInvalid}>
                        <Label htmlFor={field.name}>
                            Email
                        </Label>
                        <Input
                            id={field.name}
                            name={field.name}
                            value={field.state.value}
                            onBlur={field.handleBlur}
                            onChange={e => field.handleChange(e.target.value)}
                            aria-invalid={isInvalid}
                            placeholder="Email"
                            autoComplete="off"
                            disabled={sendLink.isPending}
                        />
                        <FieldDescription>
                            We&apos;ll email you a login code and link.
                        </FieldDescription>
                        {isInvalid && (
                            <FieldError errors={field.state.meta.errors} />
                        )}
                    </Field>
                )
            }}
        </FromField>
        <Field>
            <Button type="submit" disabled={sendLink.isPending}>
                {sendLink.isPending ? "Sending link…" : "Login"}
            </Button>
            {sendLink.isError && (
                <FieldError>{sendLink.error.message}</FieldError>
            )}
        </Field>
    </form>
}

type LoginCodeFormProps = {
    email: string
    onUseDifferentEmail: () => void
}

// TODO: Move this
function LoginCodeForm({ email, onUseDifferentEmail }: LoginCodeFormProps) {
    const verify = useMutation({
        mutationFn: (code: string) => verifyLoginCode(email, code),
        // Full navigation: a soft nav keeps the intercepted @login modal slot
        // mounted over /dashboard.
        // eslint-disable-next-line @next/next/no-location-assign-relative-destination
        onSuccess: () => window.location.assign("/dashboard")
    })

    const { handleSubmit, Field: FromField } = useForm({
        defaultValues: {
            code: ""
        },
        validators: {
            onSubmit: codeSchema
        },
        onSubmit: ({ value }) => {
            verify.mutate(value.code)
        }
    })

    const isBusy = verify.isPending || verify.isSuccess

    return <form
        className="flex flex-col gap-4 w-full"
        id="login-code-form"
        onSubmit={(e) => {
            e.preventDefault()
            handleSubmit()
        }}
    >
        <p className="text-sm">
            We sent a code and a login link to {email}. Enter the code here, or click the link in the email.
        </p>
        {process.env.NODE_ENV === "development" && (
            <p className="text-sm text-muted-foreground">
                Running locally? The code and link are printed in the API terminal.
            </p>
        )}
        <FromField name="code">
            {(field) => {
                const isInvalid =
                    field.state.meta.isTouched && !field.state.meta.isValid
                return (
                    <Field data-invalid={isInvalid}>
                        <Label htmlFor={field.name}>
                            Login code
                        </Label>
                        <Input
                            id={field.name}
                            name={field.name}
                            value={field.state.value}
                            onBlur={field.handleBlur}
                            onChange={e => field.handleChange(e.target.value.replace(/\D/g, ""))}
                            aria-invalid={isInvalid}
                            placeholder="123456"
                            inputMode="numeric"
                            autoComplete="one-time-code"
                            maxLength={6}
                            autoFocus
                            disabled={isBusy}
                        />
                        {isInvalid && (
                            <FieldError errors={field.state.meta.errors} />
                        )}
                    </Field>
                )
            }}
        </FromField>
        <Field>
            <Button type="submit" disabled={isBusy}>
                {isBusy ? "Logging in…" : "Log in"}
            </Button>
            {verify.isError && (
                <FieldError>{verify.error.message}</FieldError>
            )}
            <Button type="button" variant="ghost" onClick={onUseDifferentEmail} disabled={isBusy}>
                Use a different email
            </Button>
        </Field>
    </form>
}
