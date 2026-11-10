import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const head = el.querySelector("g");
  animate(head, { rotate: [0, -6, 4, -2, 0] }, { duration: τmax });

  const eye = el.querySelector('path[d="M15.5 10H15"]');
  animate(eye, { scaleY: [1, 0.1, 1] }, { duration: τmin, delay: τmid });
}
