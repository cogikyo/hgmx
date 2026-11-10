import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const balloon = el.querySelector("g");
  animate(balloon, { translateY: [0, -3, 1, 0] }, { duration: τmax });

  const string = el.querySelector('path[d^="M12 16"]');
  animate(string, { rotate: [0, 12, -6, 0] }, { duration: τmid, delay: τmin });
}
