"use client"

import { LoginForm } from "@/components/LoginForm/LoginForm";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from "@/components/ui/dialog";
import { useRouter } from "next/navigation";

export default function Page() {
    const { back } = useRouter()
    return <Dialog open onOpenChange={back}>
        <DialogContent>
            <DialogHeader>
                <DialogTitle>Login with your email</DialogTitle>
                <DialogDescription>
                    No password needed. We&apos;ll email you a code to enter here, or a link you can click instead.
                </DialogDescription>
            </DialogHeader>
            <LoginForm />
        </DialogContent>
    </Dialog>
}