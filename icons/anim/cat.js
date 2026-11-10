import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const cat = el.querySelector("g");
  animate(cat, { scale: [1, 0.9, 1.05, 1] }, { duration: τmax });

  const eyes = el.querySelectorAll('path[d$="v.5"]');
  animate(eyes, { scaleY: [1, 0.1, 1] }, { duration: τmin, delay: τmid });
}
