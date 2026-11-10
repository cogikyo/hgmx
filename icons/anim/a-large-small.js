import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const large = el.querySelectorAll('path[d^="m2 16"], path[d^="M3.304"]');
  animate(large, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  const small = el.querySelector("g");
  animate(small, { scale: [1, 1.25, 1] }, { duration: τmid, delay: τmin });
}
