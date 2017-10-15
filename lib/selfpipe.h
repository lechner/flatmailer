#ifndef FLATMAILER_SELFPIPE__H__
#define FLATMAILER_SELFPIPE__H__

class selfpipe
{
 public:
  selfpipe();

  operator bool() const;

  void catchsig(int sig);
  int caught();
  int waitsig(int timeout = 0);
};

#endif // FLATMAILER_SELFPIPE__H__
