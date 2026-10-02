"use client"

import { me } from "@/lib/api"
import { logout, type User } from "@/lib/auth"
import { useMutation, useQuery, useQueryClient, type UseMutationResult } from "@tanstack/react-query"
import { useRouter } from "next/navigation"
import { createContext, useContext, type ReactNode } from "react"

type CurrentUserContextType = {
    user: User | null,
    isLoading: boolean,
    logout: UseMutationResult<void, Error, void>
}

const CurrentUserContext = createContext<CurrentUserContextType | null>(null)

export function CurrentUserProvider({ children }: { children: ReactNode }) {
    const router = useRouter()
    const queryClient = useQueryClient()

    const { data: user = null, isLoading } = useQuery({
        queryKey: ["user"],
        queryFn: me,
        staleTime: 5 * 60_000
    })

    const logoutCurrentUser = useMutation({
        mutationFn: logout,
        onSuccess: () => {
            queryClient.setQueryData(["user"], null)
            queryClient.removeQueries({ predicate: query => query.queryKey[0] !== "user" })
            router.replace("/")
        }
    })

    return <CurrentUserContext.Provider value={{
        user,
        isLoading,
        logout: logoutCurrentUser
    }}>
        {children}
    </CurrentUserContext.Provider>
}

export function useCurrentUser() {
    const context = useContext(CurrentUserContext)
    if (!context) {
        throw new Error("useCurrentUser must be used inside a CurrentUserProvider")
    }
    return context
}