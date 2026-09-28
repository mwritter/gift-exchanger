"use client"

import { Button } from "@/components/ui/button";
import { Card, CardAction, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Calendar, Edit, Trash, User2, Wallet } from "lucide-react";
import type { Exchange } from "@giftexchanger/types";
import { parseISO } from "date-fns";
import Link from "next/link";

export function ExchangeListItems({ id, name, exchangeDate, budgetCents, description, memberCount, organizerId }: Exchange) {

    const handleDeleteExchange = () => {
        console.log("deleting exchange " + id)
    }

    // Get the current user id

    const isOrganizer = organizerId === 'current-user-id'

    // if not isOrganizer, fetch the organizer user to get the display name and email

    return <li>
        <Card>
            <CardHeader>
                <CardTitle>{name}</CardTitle>
                <CardDescription>{description}</CardDescription>
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
                {memberCount > 1 && <div className="flex items-center gap-2">
                    <User2 size={15} /> {memberCount} participants
                </div>}
                <div className="flex items-center gap-2">
                    <Calendar size={15} /> {parseISO(exchangeDate).toDateString()}
                </div>
                {budgetCents != null && <div className="flex items-center gap-2">
                    <Wallet size={15} /> ${budgetCents / 100} budget
                </div>}
            </CardContent>
            <CardFooter>
                {/* Show organizer display name || email */}
                <p className="text-sm">Created by you</p>
            </CardFooter>
        </Card>
    </li>
}