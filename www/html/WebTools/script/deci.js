/**
 * Created by waz on 2016-09-27.
 */

/**
 *
 *  
 *  用于加密字符串解密
 *
 **/


function EncodeString(msg) {

	return ChangeString(ChangeString(BASE64.encoder(msg)));
}

function DecodeString(msg) {
	var str = "" ;
	var unicode= BASE64.decoder(ChangeString1(ChangeString1(msg)));
	for(var i = 0 , len =  unicode.length ; i < len ;++i){
		str += String.fromCharCode(unicode[i]); 
	}
	return str ;
}


function ChangeString(str) {
	var rst = "" ;
	var aStr = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789=+/";
	var bStr = "OrsLPFGwQCHzASR67efghYUd5VWKy+/90=pqM48X2ZaNbcT3D1IJvijklmnotuBEx";
	

	for( var i = 0 ; i<str.length ; i++) {
		var ind = aStr.indexOf(str.charAt(i));
		if( ind > -1 ) {
			rst += bStr.charAt(ind) ;
		} else {
			rst += str.charAt(i) ;
		}
		
	}
	
	return rst ;
}


function ChangeString1(str) {
	var rst = "" ;
	var aStr = "OrsLPFGwQCHzASR67efghYUd5VWKy+/90=pqM48X2ZaNbcT3D1IJvijklmnotuBEx";
	var bStr = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789=+/";
	

	for( var i = 0 ; i<str.length ; i++) {
		var ind = aStr.indexOf(str.charAt(i));
		if( ind > -1 ) {
			rst += bStr.charAt(ind) ;
		} else {
			rst += str.charAt(i) ;
		}
		
	}
	
	return rst ;
}
