import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bubble = el.querySelector('path[d^="M12.7"]');
  animate(bubble, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const dot = el.querySelector("circle");
  animate(dot, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
