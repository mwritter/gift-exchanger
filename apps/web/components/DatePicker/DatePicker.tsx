"use client"

import * as React from "react"
import { format } from "date-fns"

import { Button } from "@/components/ui/button"
import { Calendar } from "@/components/ui/calendar"
import { Field, FieldError, FieldLabel } from "@/components/ui/field"
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover"

type Props = {
    title: string
    id?: string
    value: Date | null
    onChange: (date: Date | null) => void
    onBlur?: () => void
    invalid?: boolean
    errors?: Array<{ message?: string } | undefined>
}

export function DatePicker({
    id,
    title,
    value,
    onChange,
    onBlur,
    invalid,
    errors,
}: Props) {
    const picker = id ?? title

    return (
        <Field className="w-full" data-invalid={invalid}>
            <FieldLabel htmlFor={`date-picker-${picker}`}>{title ?? "Date"}</FieldLabel>
            <Popover>
                <PopoverTrigger
                    render={
                        <Button
                            variant="outline"
                            id={`date-picker-${picker}`}
                            className="justify-start font-normal"
                            aria-invalid={invalid}
                            onBlur={onBlur}
                        >
                            {value ? format(value, "PPP") : <span className="text-muted-foreground">December 24th, {(new Date()).getFullYear()}</span>}
                        </Button>
                    }
                />
                <PopoverContent className="w-auto p-0" align="start">
                    <Calendar
                        mode="single"
                        selected={value ?? undefined}
                        onSelect={(date) => onChange(date ?? null)}
                        defaultMonth={value ?? undefined}
                    />
                </PopoverContent>
            </Popover>
            {invalid && <FieldError errors={errors} />}
        </Field>
    )
}
