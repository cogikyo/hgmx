import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bubble = el.querySelector('path[d^="M22 17"]');
  animate(bubble, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const heart = el.querySelector('path[d^="M7.5 9.5"]');
  animate(heart, { scale: [1, 1.25, 0.9, 1.15, 1] }, { duration: τmid, delay: τmin });
}
