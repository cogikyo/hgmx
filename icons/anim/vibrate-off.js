import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const phone = el.querySelectorAll('path[d^="M8 8"], path[d^="M16 10.34"]');
  animate(phone, { translateX: [0, -1, 1, 0] }, { duration: τmid });

  const slash = el.querySelector("line");
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
