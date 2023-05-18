#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <pthread.h>
#include <time.h>

#define TEST_LEN 1000
#define NUM_WORKERS 10
int input[2];
int output[2];
struct gcdmsg {
    int a;
    int b;
    int out;

};

void gcd(){
       //     printf("DOING GCD STUFF!!\n");
        while(1){
        struct gcdmsg gcdmsg;
        read(input[0], &gcdmsg,sizeof(struct gcdmsg));
        int a = gcdmsg.a;
        int b = gcdmsg.b;
        printf("Calculating GCD Between %d and %d\n", gcdmsg.a, gcdmsg.b);
        int i = 0;
        if (a>b) {
			i = a;
		}
        else{
			i=b;
		}
        int gcd = 0;
        for (gcd = i; gcd >1;gcd--){
            if (a%gcd == 0 && b%gcd == 0){
                gcdmsg.out = gcd;
                break;
            }
        }
        if (gcd == 1){
            gcdmsg.out = 1;
        }
       //  printf("Calculated GCD Between %d and %d\n", gcdmsg.a, gcdmsg.b);
        write(output[1], &gcdmsg,sizeof(struct gcdmsg));

        }


}


int main(){
    double elapsed;
    struct timespec start, finish;
    if (pipe(input) < 0)
        exit(1);
    if (pipe(output) < 0)
        exit(1);

    pthread_t p[NUM_WORKERS];
    struct gcdmsg gcdchannel[TEST_LEN];
        printf("HERE1\n");
    for (int k = 0; k < NUM_WORKERS; k++) {
        /* read pipe */
        pthread_create(&p[k],NULL,gcd,NULL);
    }  
     clock_gettime(CLOCK_MONOTONIC, &start);
     for (int k = 0;k<TEST_LEN;k++ ){
         gcdchannel[k].a = rand() % 1000000;
         gcdchannel[k].b = rand() % 1000000;
         printf("Sending %d and %d\n", gcdchannel[k].a, gcdchannel[k].b);
         write(input[1], &gcdchannel[k],sizeof(struct gcdmsg));
     }

    for (int k = 0;k<TEST_LEN;k++ ){
        struct gcdmsg gcdOut;
        read(output[0], &gcdOut,sizeof(struct gcdmsg));
        printf("GCD between %d and %d is %d\n",gcdOut.a, gcdOut.b,gcdOut.out);
    }
    clock_gettime(CLOCK_MONOTONIC, &finish);
      elapsed = finish.tv_sec - start.tv_sec;
  elapsed += (finish.tv_nsec - start.tv_nsec) / 1000000000.0;
    printf("Duration of task is %.12f seconds", elapsed);

}