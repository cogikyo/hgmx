import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const church = el.querySelectorAll('path[d^="M6 21V7"], path[d^="m18 9"]');
  animate(church, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const cross = el.querySelector("g");
  animate(cross, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
