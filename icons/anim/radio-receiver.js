import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector("rect");
  animate(body, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const dial = el.querySelector('path[d="M18 12h.01"]');
  animate(dial, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: τmin });
}
