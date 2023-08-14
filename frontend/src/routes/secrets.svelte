<script lang="ts">
    import {auth} from '$lib/auth'
    import backendSecretsUrl from '$lib/config';
    import type {UserCredential} from 'firebase/auth';
    import {Button} from 'flowbite-svelte';
    import KayHeader from "../components/KayHeader.svelte";
    import type {Secret} from "$lib/secret";
    import AddSecretModal from "../components/AddSecretModal.svelte";
    import SearchAddBar from "../components/SearchAddBar.svelte";
    import SecretsList from "../components/SecretsList.svelte";

    export let currentUser: UserCredential | null
    let filterValue = "", secrets: Secret[] = [], openModal = false;

    $: shown = secrets.filter((s: Secret) => {
        let trimmed = filterValue.trim()
        const regex = /[A-Z]/
        let hasOnlyLower = trimmed.match(regex) == null
        if (hasOnlyLower) {
            return s.value.toLowerCase().indexOf(trimmed) != -1
        } else {
            return s.value.indexOf(trimmed) != -1
        }
    })

    async function getSecretsFromServer() {
        let user = currentUser as UserCredential
        return user.user.getIdToken().then(async token => {
            console.log(token)
            return fetch(backendSecretsUrl, {
                method: "GET",
                headers: {
                    "Authorization": "Bearer" + token
                }
            })
                .then(resp => resp.json())
                .then(body => {
                    secrets = body as Secret[]
                    console.log("got secrets", secrets)
                })
        });
    }

    async function uploadSecret(s: Secret) {
        let user = currentUser as UserCredential
        return user.user.getIdToken().then(async token => {
            console.log(token)
            fetch(backendSecretsUrl, {
                method: "POST",
                mode: "cors",
                cache: "no-cache",
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": "Bearer" + token
                },
                body: JSON.stringify({
                    value: s.value,
                    key: s.key,
                    url: s.url
                })
            })
                .then(resp => resp.json())
                .then(body => {
                    secrets = body
                    console.log(secrets)
                })
        })
    }

    function logout() {
        currentUser = null
        auth.signOut()
    }

    getSecretsFromServer()
</script>

<KayHeader text="Keeping " />
<div class="p-3">
    <AddSecretModal openModal={openModal} uploadSecret={uploadSecret}/>

    <div class="gap-4 items-center flex flex-col">
        <Button class="lg:top-4 lg:right-4" on:click={logout}>Logout</Button>
    </div>

    <br/>

    <SearchAddBar bind:clickedAdd={openModal} bind:filterValue={filterValue}/>

    <br/>

    <div class="flex justify-center">
        <SecretsList secrets={shown}/>
    </div>
</div>