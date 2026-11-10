import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const torch = el.querySelector("g");
  animate(torch, { rotate: [0, -8, 4, 0] }, { duration: τmid });

  const slash = el.querySelector('path[d="m2 2 20 20"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
