"use client"

import { useCurrentUser } from "@/components/CurrentUser/CurrentUserProvider";
import { DeleteExchangeButton } from "@/components/DeleteExchangeButton/DeleteExchangeButton";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Calendar, Edit, User2, Wallet } from "lucide-react";
import type { Exchange } from "@giftexchanger/types";
import { parseISO } from "date-fns";
import Link from "next/link";

export function ExchangeListItems({ id, name, exchangeDate, budgetCents, description, members, organizerId }: Exchange) {
    const { user } = useCurrentUser()
    const isOrganizer = organizerId === user?.id

    const organizer = members.find(member => member.isOrganizer)
    const membersText = members.length > 1 ? 'members' : 'member'

    return <li>
        <Card>
            <CardHeader>
                <CardTitle>
                    <Link href={`/dashboard/exchanges/${id}`} className="hover:underline">{name}</Link>
                </CardTitle>
                <CardDescription>{description}</CardDescription>
            </CardHeader>
            <CardContent>
                {members.length > 0 && <div className="flex items-center gap-2">
                    <User2 size={15} /> <span className="truncate">{members.length} {membersText} </span>
                </div>}
                <div className="flex items-center gap-2">
                    <Calendar size={15} /> {parseISO(exchangeDate).toDateString()}
                </div>
                {budgetCents != null && <div className="flex items-center gap-2">
                    <Wallet size={15} /> ${budgetCents / 100} budget
                </div>}
            </CardContent>
            <CardFooter className="flex justify-between">
                <p className="text-sm">Created by {isOrganizer ? "you" : organizer?.displayName || organizer?.email}</p>
                {isOrganizer && <div className="flex gap-2">
                    <Button nativeButton={false} render={<Link href={`/dashboard/exchanges/${id}/edit`} />} variant="outline">
                        <Edit size={12} />
                    </Button>
                    <DeleteExchangeButton exchangeId={id} exchangeName={name} />
                </div>}
            </CardFooter>
        </Card>
    </li>
}
