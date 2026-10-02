"use client"

import { useForm } from '@tanstack/react-form'
import { Field, FieldError } from "@/components/ui/field"

import z from 'zod'
import { Label } from '../ui/label'
import { Input } from '../ui/input'
import { Button } from '../ui/button'
import { DatePicker } from '../DatePicker/DatePicker'
import { Textarea } from '../ui/textarea'
import { ExchangeInviteForm } from '../ExchangeInviteInput/ExchangeInviteForm'

const formSchema = z.object({
    name: z.string().min(2),
    exchangeDate: z
        .union([z.date(), z.null()])
        .refine((date): date is Date => date instanceof Date, {
            error: "Pick an exchange date",
        }),
    description: z.string().optional(),
    // In dollars; converted to cents when sent to the API.
    budget: z.number().min(0).nullable(),
    inviteEmails: z.email().array()
})

export type ExchangeValues = z.input<typeof formSchema>

type Props = Partial<ExchangeValues> & {
    onSubmit: (data: ExchangeValues) => void
    isLoading?: boolean
}

export function ExchangeForm(props: Props) {
    const defaultValues: ExchangeValues = {
        name: props.name ?? "",
        exchangeDate: props.exchangeDate ?? null,
        description: props.description ?? "",
        budget: props.budget ?? null,
        inviteEmails: props.inviteEmails ?? []
    }

    const { handleSubmit, Field: FromField } = useForm({
        defaultValues,
        validators: {
            onSubmit: formSchema
        },
        onSubmit: ({ value }) => {
            props.onSubmit(value)
        },
    })

    return <div className='flex flex-col'>
        <form id='-exchange-form' onSubmit={(e) => {
            e.preventDefault()
            handleSubmit()
        }}>
            <div className='flex flex-col gap-5 mb-10'>
                <FromField name="name">
                    {(field) => {
                        const isInvalid =
                            field.state.meta.isTouched && !field.state.meta.isValid
                        return (
                            <Field data-invalid={isInvalid}>
                                <Label htmlFor={field.name}>
                                    Exchange Name
                                </Label>
                                <Input
                                    id={field.name}
                                    name={field.name}
                                    value={field.state.value}
                                    onBlur={field.handleBlur}
                                    onChange={e => field.handleChange(e.target.value)}
                                    aria-invalid={isInvalid}
                                    placeholder={`Christmas ${(new Date().getFullYear())}`}
                                    autoComplete="off"
                                />
                                {isInvalid && (
                                    <FieldError errors={[{ message: "Give your exchange a name" }]} />
                                )}
                            </Field>
                        )
                    }}
                </FromField>
                <FromField name="exchangeDate">
                    {(field) => {
                        const isInvalid = field.state.meta.errors.length > 0
                        return (
                            <DatePicker
                                title="Exchange Day"
                                id={field.name}
                                value={field.state.value}
                                onChange={field.handleChange}
                                onBlur={field.handleBlur}
                                invalid={isInvalid}
                                errors={field.state.meta.errors}
                            />
                        )
                    }}
                </FromField>
                <FromField name='description'>
                    {(field) => {
                        return <Field>
                            <Label htmlFor={field.name}>Exchange Description</Label>
                            <Textarea
                                id={field.name}
                                name={field.name}
                                value={field.state.value ?? ""}
                                onBlur={field.handleBlur}
                                onChange={e => field.handleChange(e.target.value)}
                                placeholder="Optional description for your exchange"
                            />
                        </Field>
                    }}
                </FromField>
                <FromField name='budget'>
                    {(field) => {
                        return <Field>
                            <Label htmlFor={field.name}>Exchange Budget</Label>
                            <div className='flex items-center'>
                                $<Input
                                    id={field.name}
                                    name={field.name}
                                    type='number'
                                    value={field.state.value ?? ""}
                                    onBlur={field.handleBlur}
                                    onChange={e => {
                                        const raw = e.target.value
                                        field.handleChange(raw === "" ? null : Number(raw))
                                    }}
                                    placeholder="50"
                                />
                            </div>
                        </Field>
                    }}
                </FromField>
                <FromField name='inviteEmails'>
                    {(field) => (
                        <Field>
                            <Label>Invite Emails</Label>
                            <ExchangeInviteForm
                                emails={field.state.value}
                                onChange={field.handleChange}
                            />
                        </Field>
                    )}
                </FromField>
            </div>
            <Field>
                <Button type="submit" disabled={props.isLoading}>
                    {props.isLoading ? 'Submitting' : 'Submit'}
                </Button>
            </Field>
        </form>
    </div>
}