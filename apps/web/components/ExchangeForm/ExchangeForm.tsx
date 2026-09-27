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
    exchangeDescription: z.string().optional(),
    exchangeBudget: z.number().nullable(),
    inviteEmails: z.email().array()
})

type ExchangeValues = z.input<typeof formSchema>

type Props = Partial<ExchangeValues> & {
    title: string,
    onSubmit: (data: ExchangeValues) => void
}

export function ExchangeForm(props: Props) {
    const defaultValues: ExchangeValues = {
        name: props.name ?? "",
        exchangeDate: props.exchangeDate ?? null,
        exchangeDescription: props.exchangeDescription ?? "",
        exchangeBudget: props.exchangeBudget ?? null,
        inviteEmails: props.inviteEmails ?? []
    }

    const { handleSubmit, Field: FromField } = useForm({
        defaultValues,
        validators: {
            onSubmit: formSchema
        },
        onSubmit: ({ value }) => {
            console.log({ value })
            console.log("Submitting")
        }
    })

    return <div className='flex flex-col'>
        <h1 className='text-2xl font-bold mb-10'>{props.title} Exchange</h1>
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
                <FromField name='exchangeDescription'>
                    {() => {
                        return <Field>
                            <Label htmlFor="textarea-description">Exchange Description</Label>
                            <Textarea id="textarea-description" placeholder="Optional description for your exchange" />
                        </Field>
                    }}
                </FromField>
                <FromField name='exchangeBudget'>
                    {() => {
                        return <Field>
                            <Label htmlFor="input-budget">Exchange Budget</Label>
                            <div className='flex items-center'>
                                $<Input id="input-budget" type='number' placeholder="50" />
                            </div>
                        </Field>
                    }}
                </FromField>
                <FromField name='inviteEmails'>
                    {() => (
                        <Field>
                            <Label>Invite Emails</Label>
                            <ExchangeInviteForm />
                        </Field>
                    )}
                </FromField>
            </div>
            <Field>
                <Button type="submit" >
                    {props.title} exchange
                </Button>
            </Field>
        </form>
    </div>
}