import { animate, τmax } from "../motion.js";

export default function (el) {
  const spool = el.querySelector("g");
  animate(spool, { rotateY: [0, 360] }, { duration: τmax });
}
