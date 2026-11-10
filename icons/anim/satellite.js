import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const craft = el.querySelector("g");
  animate(craft, { rotate: [0, -10, 6, 0] }, { duration: τmax });

  const signal = el.querySelector('path[d^="M9 21"]');
  animate(signal, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
