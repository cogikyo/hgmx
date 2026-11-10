import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const kettle = el.querySelector("g");
  animate(kettle, { rotate: [0, 10, -3, 0] }, { duration: τmid });

  const gauge = el.querySelector('path[d="M11 7v7"]');
  animate(gauge, { scaleY: [0, 1] }, { duration: τmid, delay: τmin });
}
