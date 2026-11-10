import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const moped = el.querySelectorAll("path, circle");
  animate(moped, { translateX: [0, 3, -1, 0] }, { duration: τmax });

  const body = el.querySelectorAll("path");
  animate(body, { translateY: [0, 1, 0] }, { duration: τmin, delay: τmid });
}
