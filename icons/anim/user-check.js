import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector('path[d*="21v-2"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const checkmark = el.querySelector('path[d*="17 11"]');
  animate(checkmark, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  animate(checkmark, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });

  const head = el.querySelector("circle");
  animate(head, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
