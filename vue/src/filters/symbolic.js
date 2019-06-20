const modestr = "rwxrwxrwx";

const symbolic = function(mode) {
    let pot = 1;
    let ret = "";
    var i = modestr.length;
    while (i--) {
        
      let ch = (mode & pot) != 0 ? modestr.charAt(i) : "-";
      ret = ch + ret;
      pot = pot * 2;
    }
    return ret;
};

export default symbolic;

