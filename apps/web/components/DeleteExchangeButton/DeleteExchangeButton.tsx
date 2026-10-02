"use client"

import { Button } from "@/components/ui/button"
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { deleteExchange } from "@/lib/exchanges"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { Trash } from "lucide-react"
import { useState } from "react"

type Props = {
    exchangeId: string
    exchangeName: string
    onDeleted?: () => void
}

export function DeleteExchangeButton({ exchangeId, exchangeName, onDeleted }: Props) {
    const [open, setOpen] = useState(false)
    const queryClient = useQueryClient()

    const deleteMutation = useMutation({
        mutationFn: () => deleteExchange(exchangeId),
        onSuccess: () => {
            setOpen(false)
            queryClient.invalidateQueries({ queryKey: ["exchanges"] })
            onDeleted?.()
        },
    })

    return (
        <AlertDialog open={open} onOpenChange={setOpen}>
            <AlertDialogTrigger
                render={<Button variant="destructive" aria-label={`Delete ${exchangeName}`}>
                    <Trash size={12} />
                </Button>}
            />
            <AlertDialogContent size="sm">
                <AlertDialogHeader>
                    <AlertDialogTitle>Delete {exchangeName}?</AlertDialogTitle>
                    <AlertDialogDescription>
                        This removes the exchange for everyone in it and can&apos;t be undone.
                    </AlertDialogDescription>
                    {deleteMutation.isError && (
                        <p className="text-sm text-destructive">{deleteMutation.error.message}</p>
                    )}
                </AlertDialogHeader>
                <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <AlertDialogAction
                        variant="destructive"
                        disabled={deleteMutation.isPending}
                        onClick={() => deleteMutation.mutate()}
                    >{deleteMutation.isPending ? "Deleting…" : "Delete"}</AlertDialogAction>
                </AlertDialogFooter>
            </AlertDialogContent>
        </AlertDialog>
    )
}
