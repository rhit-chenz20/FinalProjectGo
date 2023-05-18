#include <time.h>
#include <stdio.h>
#include <unistd.h>

#define TEST_LEN 250
int gcd(int a, int b){
       printf("Calculating GCD Between %d and %d\n", a, b);
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
                return gcd;
                break;
            }
        }
        return 1;

}


int main(){
    int a = 1;
    int b = 1;
    double elapsed;
    struct timespec start, finish;
     clock_gettime(CLOCK_MONOTONIC, &start);
     for (int k = 0;k<TEST_LEN;k++ ){
        a = rand() % 1000000;
        b = rand() % 1000000;
        int result = gcd(a,b);
        printf("GCD between %d and %d is %d\n",a, b, result);
     }
       clock_gettime(CLOCK_MONOTONIC, &finish);
        elapsed = finish.tv_sec - start.tv_sec;
        elapsed += (finish.tv_nsec - start.tv_nsec) / 1000000000.0;
        printf("Duration of task is %.12f seconds", elapsed);
}