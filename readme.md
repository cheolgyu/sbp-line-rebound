# 반등 구하기


## 테이블
    hist.rebound   
        파티션 테이블
        코드의 가격(저가,고가,시가,종가)별 그래프의 반등시작,종료 두개의 점을 추상화함.
    hist.create_table_rebound()   파티션 테이블 생성
        고정값: 1956 ~ 2031 (hist.price 와 동일)
    public.tb_rebound
        코드별 가격들의 반등정보
        public.view_*에서 반등테이블 join 시 사용하는 테이블

### 가격정보
    시가,저가,고가,종가

### hist.rebound 규칙 
- 이전 줄의 x2는 다음줄의 x1값 이여야됨.