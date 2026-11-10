import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const button = el.querySelector('path[d="M12 18h.01"]');
  animate(button, { scale: [1, 1.75, 1] }, { duration: τmin });

  const jar = el.querySelector("g");
  animate(jar, { rotate: [0, -5, 4, -2, 0] }, { duration: τmid, delay: τmin });
}
