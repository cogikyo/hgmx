import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const wheel = el.querySelector('path[d^="M12 7.318"]');
  animate(wheel, { translateY: [0, 3, -1, 0] }, { duration: τmid });

  const button = el.querySelector("circle");
  animate(button, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
