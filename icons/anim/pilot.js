import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const pilot = el.querySelectorAll("path, circle");
  animate(pilot, { translateY: [0, -2, 0.5, 0] }, { duration: τmax });

  const eyes = el.querySelectorAll('path[d$="v.5"]');
  animate(eyes, { opacity: [1, 0, 1] }, { duration: τmin, delay: τmid });
}
