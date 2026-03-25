import { initializeApp } from 'firebase/app';
import { getAuth, GoogleAuthProvider, signInWithPopup, signInWithEmailAndPassword, createUserWithEmailAndPassword, signOut, onAuthStateChanged, type User as FirebaseUser } from 'firebase/auth';

const firebaseConfig = {
    apiKey: import.meta.env.VITE_FIREBASE_API_KEY,
    authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN,
    projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID,
    storageBucket: import.meta.env.VITE_FIREBASE_STORAGE_BUCKET,
    messagingSenderId: import.meta.env.VITE_FIREBASE_MESSAGING_SENDER_ID,
    appId: import.meta.env.VITE_FIREBASE_APP_ID,
};

const app = initializeApp(firebaseConfig);
const auth = getAuth(app);
const googleProvider = new GoogleAuthProvider();

export const firebaseAuth = {
    signInWithGoogle: async (): Promise<FirebaseUser> => {
        const result = await signInWithPopup(auth, googleProvider);
        return result.user;
    },

    signInWithEmail: async (email: string, password: string): Promise<FirebaseUser> => {
        const result = await signInWithEmailAndPassword(auth, email, password);
        return result.user;
    },

    signUpWithEmail: async (email: string, password: string): Promise<FirebaseUser> => {
        const result = await createUserWithEmailAndPassword(auth, email, password);
        return result.user;
    },

    signOut: async (): Promise<void> => {
        await signOut(auth);
    },

    onAuthStateChanged: (callback: (user: FirebaseUser | null) => void) => {
        return onAuthStateChanged(auth, callback);
    },

    getCurrentUser: (): FirebaseUser | null => {
        return auth.currentUser;
    },

    getIdToken: async (): Promise<string | null> => {
        const user = auth.currentUser;
        if (user) {
            return await user.getIdToken();
        }
        return null;
    },
};

export { auth };
export default firebaseAuth;
