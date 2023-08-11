<script lang="ts"> 
    import { createUserWithEmailAndPassword, signInWithEmailAndPassword, GoogleAuthProvider, signInWithPopup, type UserCredential } from "firebase/auth";
    import { auth } from "$lib/auth"
    import { Img, Checkbox, Input, Label, Helper, Heading, Mark, Span, P, Button} from "flowbite-svelte";

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
        <Heading class="text-center mb-2 text-7xl" tag="h1"><Span class="font-mono" gradient>Kay</Span>Vault</Heading>
        <P class="text-center">Login or register to manage your <Mark>secrets</Mark></P>

        <br>

        <div class="flex justify-center">
            <div class="w-3/5 p-6 rounded bg-gradient-to-r from-sky-500 to-emerald-600 dark:bg-gray-600">
                <form>
                    <div class="mb-3">
                        <Button color="alternative">
                            <Img src="https://img.icons8.com/color/16/000000/google-logo.png"/>Signup Using Google
                        </Button>
                    </div>
                    <div class="mb-3">
                        <Label>Email address</Label>
                        <Input type="text" bind:email placeholder="kay@vault.com"></Input>
                        <Helper>We use Firebase Authentication to store user login information.</Helper>
                    </div>
                    <div class="mb-3">
                        <Label>Password</Label>
                        <Input type="password" bind:password placeholder="much secret"></Input>
                    </div>
                    <div class="mb-3">
                        <Checkbox>Remember me</Checkbox>
                    </div>
                    <Button on:click={loginUser}>Login</Button>
                    <Button color="alternative" on:click={registerNewUser}>Register</Button>
                </form>
            </div>
        </div>

    </div>
</div>