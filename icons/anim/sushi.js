import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const roll = el.querySelector('rect[width="18"]');
  animate(roll, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const filling = el.querySelector("g");
  animate(filling, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
