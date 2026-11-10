import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bowl = el.querySelectorAll('path[d^="M4 12"], path[d^="M7 21"]');
  animate(bowl, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const leaves = el.querySelector("g");
  animate(leaves, { translateY: [0, -3, 1, 0] }, { duration: τmid, delay: τmin });
}
