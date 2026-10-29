import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const tablet = el.querySelector("g");
  animate(tablet, { rotate: [-90, 10, 0] }, { duration: τmax });

  const button = el.querySelector("line");
  animate(button, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
