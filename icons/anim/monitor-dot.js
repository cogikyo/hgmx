import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const screen = el.querySelector('path[d^="M22 12.307"]');
  animate(screen, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const dot = el.querySelector("circle");
  animate(dot, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
