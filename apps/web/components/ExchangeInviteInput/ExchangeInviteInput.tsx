import z from "zod";
import { Field, FieldError } from "../ui/field";
import { Input } from "../ui/input";
import { useForm } from "@tanstack/react-form";
import { Button } from "../ui/button";

type Props = {
    onSubmit: (email: string) => void
}


const inputSchema = z.object({
    email: z.email()
})

export function ExchangeInviteInput({ onSubmit }: Props) {

    const { handleSubmit, Field: FormField, reset } = useForm({
        defaultValues: {
            email: ""
        },
        validators: {
            onSubmit: inputSchema
        },
        onSubmit: (e) => {
            onSubmit(e.value.email)
            reset()
        }
    })

    return <div className="flex flex-col gap-4 w-full"
        id={`invite-form`}>
        <FormField name="email">
            {(field) => {
                const isInvalid =
                    field.state.meta.isTouched && !field.state.meta.isValid
                return (
                    <div>
                        <div data-invalid={isInvalid} className="flex gap-2">
                            <Input
                                id={field.name}
                                name={field.name}
                                value={field.state.value}
                                onBlur={field.handleBlur}
                                onChange={e => field.handleChange(e.target.value)}
                                aria-invalid={isInvalid}
                                placeholder="Invite Email"
                                autoComplete="off"
                                onKeyDown={e => {
                                    if (e.key === "Enter") {
                                        e.preventDefault()
                                        handleSubmit()
                                    }
                                }}
                            />
                            <Button type="button" onClick={e => {
                                e.stopPropagation()
                                handleSubmit()
                            }}>
                                Add
                            </Button>
                        </div>

                        {isInvalid && (
                            <FieldError errors={field.state.meta.errors} />
                        )}
                    </div>
                )
            }}
        </FormField>
    </div>
}