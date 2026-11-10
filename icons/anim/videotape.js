import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const shell = el.querySelectorAll('rect, path[d^="M2 8"]');
  animate(shell, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const reels = el.querySelectorAll("circle");
  animate(reels, { rotate: [0, 360], pathLength: [0.25, 1] }, { duration: τmid, delay: τmin });
}
