import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const eyes = el.querySelectorAll('circle, path[d$="7h.01"]');
  animate(eyes, { translateY: [0, -2, 0.5, 0] }, { duration: τmid });

  const mouth = el.querySelector('path[d^="M22 13.5"]');
  animate(mouth, { scaleY: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
