"use client"

import { Button } from "@/components/ui/button";
import { Card, CardAction, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Calendar, Edit, Trash, User2, Wallet } from "lucide-react";
import { Exchange } from "./ExchangeList";
import Link from "next/link";

export function ExchangeListItems({ id, name, exchangeDate, exchangeBudgent, exchangeDescription, exchangeInites, exchangeOrganizerId }: Exchange) {

    const handleDeleteExchange = () => {
        console.log("deleting exchange " + id)
    }

    // Get the current user id

    const isOrganizer = exchangeOrganizerId === 'current-user-id'

    // if not isOrganizer, fetch the organizer user to get the display name and email

    return <li>
        <Card>
            <CardHeader>
                <CardTitle>{name}</CardTitle>
                <CardDescription>{exchangeDescription}</CardDescription>
                {isOrganizer && <CardAction>
                    <Button render={<Link href={`/dashboard/exchanges/edit/${id}`} />} variant={'ghost'}>
                        <Edit size={12} />
                    </Button>
                    <Button variant={'destructive'} onClick={handleDeleteExchange}>
                        <Trash size={12} />
                    </Button>
                </CardAction>}
            </CardHeader>
            <CardContent>
                {exchangeInites.length > 1 && <div className="flex items-center gap-2">
                    <User2 size={15} /> {exchangeInites.length} participants
                </div>}
                <div className="flex items-center gap-2">
                    <Calendar size={15} /> {exchangeDate.toDateString()}
                </div>
                {exchangeBudgent && <div className="flex items-center gap-2">
                    <Wallet size={15} /> ${exchangeBudgent} budget
                </div>}
            </CardContent>
            <CardFooter>
                {/* Show organizer display name || email */}
                <p className="text-sm">Created by you</p>
            </CardFooter>
        </Card>
    </li>
}