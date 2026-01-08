// Compile with `gcc foo.c -Wall -std=gnu99 -lpthread`, or use the makefile
// The executable will be named `foo` if you use the makefile, or `a.out` if you use gcc directly

#include <pthread.h>
#include <stdio.h>

int i = 0;
pthread_mutex_t i_lock; // Mutex to protect access to i, Using a mutex because i is shared between threads and not a semaphore because we protect a variable not a resource
// Note the return type: void*
void* incrementingThreadFunction(){
    // TODO: increment i 1_000_000 times
    for(int j = 0; j < 1000000; j++){
        pthread_mutex_lock(&i_lock);
        i++;
        pthread_mutex_unlock(&i_lock);
    }
    return NULL;
}

void* decrementingThreadFunction(){
    // TODO: decrement i 1_000_000 times
    for(int j = 0; j < 100000; j++){
        pthread_mutex_lock(&i_lock);
        i--;
        pthread_mutex_unlock(&i_lock);
    }   
    return NULL;
}


int main(){
    // TODO: 
    // start the two functions as their own threads using `pthread_create`
    // Hint: search the web! Maybe try "pthread_create example"?
    pthread_t incrementingThread;
    pthread_t decrementingThread;

    pthread_mutex_init(&i_lock, NULL);
    pthread_create(&incrementingThread, NULL, incrementingThreadFunction, NULL);
    pthread_create(&decrementingThread, NULL, decrementingThreadFunction, NULL);
    
    // TODO:
    // wait for the two threads to be done before printing the final result
    // Hint: Use `pthread_join`    
    pthread_join(incrementingThread, NULL);
    pthread_join(decrementingThread, NULL);

    pthread_mutex_destroy(&i_lock);
    
    printf("The magic number is: %d\n", i);
    return 0;
}