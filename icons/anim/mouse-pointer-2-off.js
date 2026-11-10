import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const pointer = el.querySelector("g");
  animate(pointer, { scale: [1, 0.92, 1] }, { duration: τmid });

  const slash = el.querySelector('path[d="M22 2 2 22"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
