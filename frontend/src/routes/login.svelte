<script lang="ts">
    import {
        createUserWithEmailAndPassword,
        GoogleAuthProvider,
        signInWithEmailAndPassword,
        signInWithPopup,
        type UserCredential
    } from "firebase/auth";
    import {auth} from "$lib/auth"
    import {Button, Checkbox, Helper, Img, Input, Label} from "flowbite-svelte";
    import {Icon} from "flowbite-svelte-icons";
    import KayHeader from "../components/KayHeader.svelte";

    let email = "";
    let password = "";
    export let currentUser: UserCredential | null = null

    function registerNewUser(): void {
        createUserWithEmailAndPassword(auth, email, password)
            .then((userCredential) => { currentUser = userCredential;
            })
            
            .catch((error) => {
                const errorCode = error.code;
                const errorMessage = error.message;
            });
    }

    function loginUser(): void {
        signInWithEmailAndPassword(auth, email, password)
            .then((userCredential) => {
                currentUser = userCredential;
            })
            .catch((error) => {

                // email and password did not match
            });
    }

    function signInWithGoogle(): void {
        const provider = new GoogleAuthProvider()
        signInWithPopup(auth, provider)
            .then((userCredential) => {
                currentUser = userCredential;
            });
    }
</script>

<div class="bg-gradient-to-b from-gray-50">
    <div class="p-3">
        <KayHeader text="Login or register to keep your"/>

        <br>

        <div class="flex justify-center">
            <div class="lg:w-1/2 2xl:w-1/4 p-6 rounded bg-gradient-to-r from-sky-500 to-emerald-600 dark:bg-gray-600">
                <form>
                    <div class="mb-3">
                        <Button color="alternative" on:click={signInWithGoogle}>
                            <Img src="https://img.icons8.com/color/16/000000/google-logo.png"/>Google Login
                        </Button>
                    </div>
                    <div class="mb-3">
                        <Label class="space-y-2">
                            <span>Email address</span>
                            <Input type="email" size="md" bind:value={email} placeholder="kay@vault.com">
                                <Icon name="envelope-solid" slot="left"></Icon>
                            </Input>
                            <Helper>We use Firebase Authentication to store user login information.</Helper>
                        </Label>
                    </div>
                    <div class="mb-3">
                        <Label class="space-y-2">
                            <span>Password</span>
                            <Input type="password" size="md" bind:value={password} placeholder="much secret"></Input>
                        </Label>
                    </div>
                    <div class="mb-3">
                        <Checkbox>Remember me</Checkbox>
                    </div>
                    <div class="grid lg:grid-cols-2 gap-4">
                        <Button on:click={loginUser}>Login</Button>
                        <Button color="alternative" on:click={registerNewUser}>Register</Button>
                    </div>
                </form>
            </div>
        </div>

    </div>
</div>