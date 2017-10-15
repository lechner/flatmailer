#ifndef FLATMAILER__MAKEFIELD__H__
#define FLATMAILER__MAKEFIELD__H__

extern mystring make_date(time_t t = 0);
extern mystring make_messageid(const mystring& idhost);
extern mystring make_boundary();

#endif // FLATMAILER__MAKEFIELD__H__
