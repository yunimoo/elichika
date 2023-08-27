# build database to /static from assets/db

python bin/encryptdbset.py static update assets/db/gl/masterdata_a_en assets/db/gl/masterdata_a_ko assets/db/gl/masterdata_a_zh assets/db/gl/masterdata_i_en assets/db/gl/masterdata_i_ko assets/db/gl/masterdata_i_zh
python bin/encryptdbset.py static update assets/db/jp/masterdata_a_ja assets/db/jp/masterdata_i_ja