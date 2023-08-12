<script lang="ts">
    import { auth } from '$lib/auth'
    import backendSecretsUrl from '$lib/config';
    import type { UserCredential } from 'firebase/auth';
    import {Button, Input} from 'flowbite-svelte';
    import { Icon } from 'flowbite-svelte-icons';
    import KayHeader from "../components/KayHeader.svelte";
    import type { Secret } from "$lib/secret";
    import AddSecretModal from "../components/AddSecretModal.svelte";

    export let currentUser: UserCredential | null
    let filterValue = "", secrets: Secret[] = [];

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

<style>
    .input {
        text-align: center;
    }
    .secrets {
        text-align: center;
    }
    .create {
        text-align: center;
        margin: auto;
    }
</style>

<html lang="en">
<body>

<div class="container p-3">
    <AddSecretModal uploadSecret={uploadSecret}/>
    <KayHeader text="Manage your"/>

    <br/>

    <div class="container text-center">
        <div class="row">
            <div class="input col-10">
                <Input bind:value={filterValue} type="text" size="lg" placeholder="Add/Search Secrets" name="New Secret">
                    <Icon name="search-outline" slot="left"></Icon>
                </Input>
            </div>
            <div class="col create">
                <Button color="blue" data-bs-toggle="modal" data-bs-target="#exampleModal">Create</Button>
            </div>
        </div>
    </div>
    <div class="secrets">
        {#each shown as secret}
            <p>Name:{secret.key} Value:{secret.value}, Url:{secret.url}</p>
        {/each}
    </div>
</div>
</body>
</html>
